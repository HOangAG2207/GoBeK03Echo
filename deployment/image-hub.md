login docker
tạo repository trên hub.docker.com
cd root project
docker build -t ten-repo-image .
docker tag <ten-repo-image> <user-hub>/<ten-repo-image>:<tag>
docker push <user-hub>/<ten-repo-image>:<tag>