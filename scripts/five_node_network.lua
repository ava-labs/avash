cmds = {
    "startnode n1 --db-enabled=false --staking-tls-enabled=false --api-ipcs-enabled=false --http-port=9652 --staking-port=9155 --xput-server-port=9255 --log-level=verbo --bootstrap-ips= --staking-tls-cert-file=certs/keys1/staker.crt --staking-tls-key-file=certs/keys1/staker.key",
    "startnode n2 --db-enabled=false --staking-tls-enabled=false --api-ipcs-enabled=false --http-port=9654 --staking-port=9156 --xput-server-port=9256 --log-level=verbo --bootstrap-ips=127.0.0.1:9155 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg --staking-tls-cert-file=certs/keys2/staker.crt --staking-tls-key-file=certs/keys2/staker.key",
    "startnode n3 --db-enabled=false --staking-tls-enabled=false --api-ipcs-enabled=false --http-port=9656 --staking-port=9157 --xput-server-port=9257 --log-level=verbo --bootstrap-ips=127.0.0.1:9155,127.0.0.1:9156 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg,MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ --staking-tls-cert-file=certs/keys3/staker.crt --staking-tls-key-file=certs/keys3/staker.key",
    "startnode n4 --db-enabled=false --staking-tls-enabled=false --api-ipcs-enabled=false --http-port=9658 --staking-port=9158 --xput-server-port=9258 --log-level=verbo --bootstrap-ips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg,MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ,NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN --staking-tls-cert-file=certs/keys4/staker.crt --staking-tls-key-file=certs/keys4/staker.key",
    "startnode apinode --db-enabled=false --staking-tls-enabled=false --api-ipcs-enabled=true --http-port=9650 --staking-port=9159 --xput-server-port=9259 --log-level=verbo --bootstrap-ips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157,127.0.0.1:9158 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg,MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ,NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN,GWPcbFJZFfZreETSoWjPimr846mXEKCtu --staking-tls-cert-file=certs/keys5/staker.crt --staking-tls-key-file=certs/keys5/staker.key",
}

for key, cmd in ipairs(cmds) do
    print("calling " .. cmd)
    avash_call(cmd)
end
