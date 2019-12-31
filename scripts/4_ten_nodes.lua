second = 1000000

nodeprefix = "a"

-- Deploys 10 nodes: a1 -> a10 
cmds = {
    "startnode " .. nodeprefix .. "1 --db=false --jrpcport=9650 --serverport=9150 --rpcport=9350 --loglevel=all --bootstrapips=",
    "startnode " .. nodeprefix .. "2 --db=false --jrpcport=9651 --serverport=9151 --rpcport=9351 --loglevel=all --bootstrapips=127.0.0.1:9150",
    "startnode " .. nodeprefix .. "3 --db=false --jrpcport=9652 --serverport=9152 --rpcport=9352 --loglevel=all --bootstrapips=127.0.0.1:9150,127.0.0.1:9151",
}

bsips = " --loglevel=all --bootstrapips=127.0.0.1:9150,127.0.0.1:9151,127.0.0.1:9152"

cmds_template = {
    " --db=false --jrpcport=965",
    "--serverport=915",
    "--rpcport=935",
}

for i=4, 10 do
    iminus1 = i - 1
    new_cmd = "startnode " .. nodeprefix .. i .. " " .. table.concat(cmds_template, iminus1 .. " ") .. iminus1
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