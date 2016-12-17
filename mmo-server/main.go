package main

import (
	"flag"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/profile"
	"github.com/sintell/mmo-server/auth"
	"github.com/sintell/mmo-server/db"
	"github.com/sintell/mmo-server/game"
	"github.com/sintell/mmo-server/packet"
	"github.com/sintell/mmo-server/server"
)

var prof = flag.String("prof", "", "write cpu profile to file")
var stopdelay = flag.Int("stopdelay", -1, "wait n secs before server stop")

var ip = flag.String("ip", "0.0.0.0", "server ip adress")
var port = flag.Int("port", 3034, "server port")
var dbIP = flag.String("dbip", "138.201.123.151", "db ip adress")
var dbPort = flag.Int("dbport", 14112, "db port")

func init() {
	flag.Parse()
}

func profl(prof string) interface {
	Stop()
} {
	switch prof {
	case "cpu":
		return profile.Start(profile.ProfilePath("./prof"), profile.NoShutdownHook, profile.CPUProfile)
	case "mem":
		return profile.Start(profile.ProfilePath("./prof"), profile.NoShutdownHook, profile.MemProfile)
	case "block":
		return profile.Start(profile.ProfilePath("./prof"), profile.NoShutdownHook, profile.BlockProfile)
	case "trace":
		return profile.Start(profile.ProfilePath("./prof"), profile.NoShutdownHook, profile.TraceProfile)
	default:
		go http.ListenAndServe(":8080", http.DefaultServeMux)
		return nil
	}
}

func main() {
	profStop := profl(*prof)

	db, err := db.NewProvider(&net.TCPAddr{IP: net.ParseIP(*dbIP), Port: *dbPort}, 1)
	if err != nil {
		glog.Fatal(err.Error())
	}

	server := server.TCPServer{
		NetAddr: &net.TCPAddr{IP: net.ParseIP(*ip), Port: *port},
		ConnectionManager: server.ConnectionManager{
			PacketHandler: &packet.GamePacketHandler{
				HeadLength: 6,
				PacketList: packet.PacketsList{},
			},
			Connections: make(map[server.TCPConnection]bool),
		},
		GameManager: game.NewManager(db),
		AuthManager: auth.NewManager(db),
		DB:          *db,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			glog.Infof("got %s, finishing...\n", sig.String())
			server.Stop()
			db.Stop()
			if *prof != "" {
				profStop.Stop()
			}

			if *stopdelay > 0 {
				glog.Infof("waiting %ss for all gorutines to finish...\n", *stopdelay)
				<-time.After(time.Second * time.Duration(*stopdelay))
			}

			glog.Flush()
			os.Exit(0)
		}
	}()

	server.Listen()
}
