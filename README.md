# mqtt-listener-go #

## snippets ##

```shell
sudo apt-get install qemu qemu-user-static binfmt-support debootstrap -y
docker buildx create --name armBuilder
docker buildx inspect --bootstrap
docker run --privileged --rm tonistiigi/binfmt --install all
```

## links ##

- [](https://www.docker.com/blog/multi-arch-images/)
- [Github buildx](https://github.com/docker/buildx/)
