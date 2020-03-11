cmds = {
    "startnode node1 --staking-tls-enabled=true --http-port=9650 --staking-port=9651 --log-level=verbo --staking-tls-enabled=false --bootstrap-ips= --staking-tls-cert-file=certs/keys1/staker.crt --staking-tls-key-file=certs/keys1/staker.key",
    "startnode node2 --staking-tls-enabled=true --http-port=9652 --staking-port=9653 --log-level=verbo --staking-tls-enabled=false --bootstrap-ips=127.0.0.1:9651 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg --staking-tls-cert-file=certs/keys2/staker.crt --staking-tls-key-file=certs/keys2/staker.key",
    "startnode node3 --staking-tls-enabled=true --http-port=9654 --staking-port=9655 --log-level=verbo --staking-tls-enabled=false --bootstrap-ips=127.0.0.1:9651 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg --staking-tls-cert-file=certs/keys3/staker.crt --staking-tls-key-file=certs/keys3/staker.key",
    "startnode node4 --staking-tls-enabled=true --http-port=9656 --staking-port=9657 --log-level=verbo --staking-tls-enabled=false --bootstrap-ips=127.0.0.1:9651 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg --staking-tls-cert-file=certs/keys4/staker.crt --staking-tls-key-file=certs/keys4/staker.key",
    "startnode node5 --staking-tls-enabled=true --http-port=9658 --staking-port=9659 --log-level=verbo --staking-tls-enabled=false --bootstrap-ips=127.0.0.1:9651 --bootstrap-ids=7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg --staking-tls-cert-file=certs/keys5/staker.crt --staking-tls-key-file=certs/keys5/staker.key",
}

for key, cmd in ipairs(cmds) do
    avash_call(cmd)
end