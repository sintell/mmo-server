# Mmo-server
This is naive approach to create a mmo-server using Go


### Comandline options:
##### Logging:
- `-logtostderr` - log program output to std instead of file
- `-alsologtostderr` - log program output to std and to file
- `-log_dir="file_path/file_name"` - log program output to `file_path/file_name`

##### Profiling:
- `-prof="type"` - launch server with profiler enabled, where `type` could be one of:
    - `cpu` - for cpu profiling
    - `mem` - for memmory profiling
    - `block` - for blocks profiling
    - `trace` - for function tracing
