REPO="github.com/j-forster/Wazihub-API-Server"

go install "$REPO"

systemctl stop wazihub
systemctl disable wazihub.service

cp "$GOPATH/bin/Wazihub-API-Server" "/bin/wazihub-gateway"
cp "$GOPATH/src/$REPO/wazihub.service" /lib/systemd/system/wazihub.service

systemctl enable wazihub.service
