 #!/bin/bash
 echo $1 | gpg \
	--passphrase-fd 0 \
	--batch --yes \
	--no-default-keyring --armor \
	--secret-keyring ./unixvoid.sec --keyring ./unixvoid.pub \
	--output go_ip_updater-latest-linux-amd64.aci.asc \
	--detach-sig go_ip_updater-latest-linux-amd64.aci
