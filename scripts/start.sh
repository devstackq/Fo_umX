app="Forum-X2"
docker run -d --rm -p 8080:6969 \
	--name=${app} \
	devstackq/${app}