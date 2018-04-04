# faiss-proxy
A http proxy for faiss-server.

# build

```shell
sh build.sh
```

# run

```shell
./faiss_proxy
```

# Restful APIs

```shell
/faiss/1.0/ping
/faiss/1.0/db/new
/faiss/1.0/db/del
/faiss/1.0/db/list
/faiss/1.0/hset
/faiss/1.0/hdel
/faiss/1.0/hget
/faiss/1.0/hsearch
```

# Example

```shell
curl -X POST -k http://localhost:3839/faiss/1.0/db/new -d '{"db_name":"1000w", "max_size":10000000}'
```

# Tips
Faiss-proxy should run with faiss-server. There's some configs for faiss-proxy, to see them just run `./faiss_proxy --help`. Parameter `-faiss_endpoint` is the endpoint of faiss-server.
