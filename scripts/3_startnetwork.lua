second = 1000000
cmds = {
    "startnode a1 --db=false --jrpcport=9655 --serverport=9155 --rpcport=9355 --loglevel=all --bootstrapips=",
    "startnode a2 --db=false --jrpcport=9656 --serverport=9156 --rpcport=9356 --loglevel=all --bootstrapips=127.0.0.1:9155",
    "startnode a3 --db=false --jrpcport=9657 --serverport=9157 --rpcport=9357 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156",
    "startnode jrpcnode --db=false --jrpcport=9650 --serverport=9150 --rpcport=9350 --loglevel=all --bootstrapips=127.0.0.1:9155,127.0.0.1:9156,127.0.0.1:9157",
}

for key, cmd in ipairs(cmds) do
    print("calling " .. cmd)
    avash_call(cmd)
end
--[[
avash_sleepmicro(1 * second)
avash_call("procmanager list")
avash_call("procmanager stopall")
avash_sleepmicro(1 * second)
avash_call("procmanager list")
avash_call("procmanager startall")
avash_sleepmicro(1 * second)
avash_call("procmanager list")
avash_call("procmanager remove a1")
avash_call("procmanager remove a2")
avash_call("procmanager remove a3")
avash_call("procmanager remove jrpcnode")
avash_sleepmicro(1 * second)
avash_call("procmanager list")
]]