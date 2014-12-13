echo "Creating docker containers"
echo "==========================="
docker ps -a
echo "==========================="
docker rm `docker ps -a -q`
docker run -d --name node1 -h node1 progrium/consul -server -bootstrap-expect 3
JOIN_IP="$(docker inspect -f '{{.NetworkSettings.IPAddress}}' node1)"
docker run -d --name node2 -h node2 progrium/consul -server -join $JOIN_IP
docker run -d --name node3 -h node3 progrium/consul -server -join $JOIN_IP
docker run -d -p 8400:8400 -p 8500:8500 -p 8600:53/udp --name node4 -h node4 progrium/consul -join $JOIN_IP
