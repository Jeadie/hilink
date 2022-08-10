DATA=`curl http://192.168.8.1/api/webserver/SesTokInfo`
SESSION_ID=`echo "$DATA" | grep "SesInfo" | cut -b 58-185`
TOKEN=`echo "$DATA" | grep "TokInfo"| cut -b 205-236`

curl http://192.168.8.1/api/sms/send-sms \
-H "Cookie: $SESSION_ID" -H "__RequestVerificationToken: $TOKEN" \
 --data "<?xml version='1.0' encoding='UTF-8'?><request><Index>-1</Index><Phones><Phone>$1</Phone></Phones><Sca></Sca><Content>$2</Content><Length>-1</Length><Reserved>1</Reserved><Date>-1</Date></request>"
