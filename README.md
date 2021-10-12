# mqtt-listener-go #

## snippets ##

```shell
sudo apt-get install qemu qemu-user-static binfmt-support debootstrap -y
docker buildx create --name armBuilder
docker buildx inspect --bootstrap
docker run --privileged --rm tonistiigi/binfmt --install all
```

### mongodb raspberry ###

```bash
echo 'deb http://ftp.br.debian.org/debian stretch main' /etc/apt/sources.list.d/repo_mongodb_org_debian.list
sudo apt-get update && sudo apt-get install mongodb-server
```

### install ###

```bash
# raspberry
bash <(curl -s -L https://raw.githubusercontent.com/Eldius/mqtt-listener-go/main/scripts/install_raspiberry.sh) --argument1=true

# amd64
bash <(curl -s -L https://raw.githubusercontent.com/Eldius/mqtt-listener-go/main/scripts/install_amd64.sh) --argument1=true
```


## links ##

- [](https://www.docker.com/blog/multi-arch-images/)
- [Github buildx](https://github.com/docker/buildx/)
