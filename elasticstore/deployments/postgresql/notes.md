While installing pgsql, in case of data directory error 

Error response from daemon: failed to populate volume: error while mounting volume '/var/lib/docker/volumes/postgresql_pg-data/_data': failed to mount local volume: mount /home/user/docker-volumes/postgres-data:/var/lib/docker/volumes/postgresql_pg-data/_data, flags: 0x1000: no such file or directory

Use below commands to create folder and give 755 persmission

mkdir -p /home/user/docker-volumes/postgres-data

chmod -R 755 /home/user/docker-volumes/postgres-data