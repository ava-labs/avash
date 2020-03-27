second = 1000000

cmds = {
    "varstore create st",
    "startnode n1 --db-enabled=false --http-port=9650 --log-level=verbo",
    "callrpc n1 ext/admin admin.getNetworkID {} st nid",
    "procmanager remove n1"
}

for key, cmd in ipairs(cmds) do
    print("calling " .. cmd)
    avash_call(cmd)
    avash_sleepmicro(1 * second)
end