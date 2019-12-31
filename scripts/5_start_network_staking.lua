second = 1000000
cmds = {
    "startnode a1 --staking_enabled=true --db=false --jrpcport=9655  --serverport=9155  --rpcport=9355 --loglevel=all --bootstrapips= --stake_cert_file=certs/keys1/staker.crt --stake_key_file=certs/keys1/staker.key",
    "startnode a2 --staking_enabled=true --db=false --jrpcport=9656  --serverport=9156  --rpcport=9356 --loglevel=all --bootstrapips=127.0.0.1:9155 --stake_cert_file=certs/keys2/staker.crt --stake_key_file=certs/keys2/staker.key",
    "startnode a3 --staking_enabled=true --db=false --jrpcport=9657  --serverport=9157  --rpcport=9357 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156 --stake_cert_file=certs/keys3/staker.crt --stake_key_file=certs/keys3/staker.key",
    "startnode a4 --staking_enabled=true --db=false --jrpcport=9658  --serverport=9158  --rpcport=9358 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157 --stake_cert_file=certs/keys4/staker.crt --stake_key_file=certs/keys4/staker.key",
    "startnode a5 --staking_enabled=true --db=false --jrpcport=9659  --serverport=9159  --rpcport=9359 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157,127.0.0.1:9158 --stake_cert_file=certs/keys5/staker.crt --stake_key_file=certs/keys5/staker.key",
}

for key, cmd in ipairs(cmds) do
    avash_call(cmd)
end