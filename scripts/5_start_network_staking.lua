second = 1000000
cmds = {
    "startnode a1 --staking-tls-enabled=true --db-enabled=false --http-port=9655 --loglevel=all --bootstrapips= --staking-tls-cert-file=certs/keys1/staker.crt --staking-tls-key-file=certs/keys1/staker.key",
    "startnode a2 --staking-tls-enabled=true --db-enabled=false --http-port=9656 --loglevel=all --bootstrapips=127.0.0.1:9155 --staking-tls-cert-file=certs/keys2/staker.crt --staking-tls-key-file=certs/keys2/staker.key",
    "startnode a3 --staking-tls-enabled=true --db-enabled=false --http-port=9657 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156 --staking-tls-cert-file=certs/keys3/staker.crt --staking-tls-key-file=certs/keys3/staker.key",
    "startnode a4 --staking-tls-enabled=true --db-enabled=false --http-port=9658 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157 --staking-tls-cert-file=certs/keys4/staker.crt --staking-tls-key-file=certs/keys4/staker.key",
    "startnode a5 --staking-tls-enabled=true --db-enabled=false --http-port=9659 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157,127.0.0.1:9158 --staking-tls-cert-file=certs/keys5/staker.crt --staking-tls-key-file=certs/keys5/staker.key",
}

for key, cmd in ipairs(cmds) do
    avash_call(cmd)
end