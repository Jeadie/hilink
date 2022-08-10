// curl 'http://192.168.8.1/api/webserver/SesTokInfo
curl 'http://192.168.8.1/api/dialup/connection' \
  -H 'Cookie: SessionID=$1' \
  -H '__RequestVerificationToken: $2'  --data '<?xml version="1.0" encoding="UTF-8"?><request><RoamAutoConnectEnable>0</RoamAutoConnectEnable><MaxIdelTime>86400</MaxIdelTime><ConnectMode>0</ConnectMode><MTU>1500</MTU><auto_dial_switch>1</auto_dial_switch><pdp_always_on>0</pdp_always_on></request>'

curl 'http://192.168.8.1/api/cradle-mac' \
  -H 'Cookie: SessionID=$1' \
  -H '__RequestVerificationToken: $2'
