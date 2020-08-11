second = 1000000

nodeprefix = "a"

-- Deploys 10 nodes: a1 -> a10 
-- Commands 1-3 deploy initial nodes for others to 
cmds = {
    "startnode " .. nodeprefix .. "1 --db-enabled=false --api-ipcs-enabled=true --http-port=9650 --staking-port=9150 --log-level=verbo --bootstrap-ips= --staking-tls-cert-file=certs/keys1/staker.crt --staking-tls-key-file=certs/keys1/staker.key",
    "startnode " .. nodeprefix .. "2 --db-enabled=false --api-ipcs-enabled=false --http-port=9651 --staking-port=9151 --log-level=verbo --bootstrap-ips=127.0.0.1:9150 --bootstrap-ids=NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg --staking-tls-cert-file=certs/keys2/staker.crt --staking-tls-key-file=certs/keys2/staker.key",
    "startnode " .. nodeprefix .. "3 --db-enabled=false --api-ipcs-enabled=false --http-port=9652 --staking-port=9152 --log-level=verbo --bootstrap-ips=127.0.0.1:9150,127.0.0.1:9151 --bootstrap-ids=NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg,NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ --staking-tls-cert-file=certs/keys3/staker.crt --staking-tls-key-file=certs/keys3/staker.key",
}

bsips = " --log-level=verbo --bootstrap-ips=127.0.0.1:9150,127.0.0.1:9151,127.0.0.1:9152 --bootstrap-ids=NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg,NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ,NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN"

cmds_template = {
    " --db-enabled=false --api-ipcs-enabled=false --http-port=965",
    "--staking-port=915",
    "--xput-server-port=925"
}

-- Add empty string into cmds_template because table.concat does not append to last element
table.insert(cmds_template, "")

for i=4, 10 do
    iminus1 = i - 1
    key_params = "--staking-tls-cert-file=certs/keys" .. i .. "/staker.crt --staking-tls-key-file=certs/keys" .. i .. "/staker.key"
    new_cmd = "startnode " .. nodeprefix .. i .. " " .. table.concat(cmds_template, iminus1 .. " ") .. key_params
    cmds[i] = new_cmd .. bsips
end

nodes = {}
for i=1, 10 do
    nodes[i] = nodeprefix .. i
end

for key, cmd in ipairs(cmds) do
    print("calling " .. cmd)
    avash_call(cmd)
    -- wait 1 second for the node to boot up
    print("sleeping 1sec")
    avash_sleepmicro(1 * second)
end