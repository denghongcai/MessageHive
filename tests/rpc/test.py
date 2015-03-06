import pyjsonrpc

http_client = pyjsonrpc.HttpClient(
    url = "http://127.0.0.1:9999/rpc"
)

print http_client.call("GroupManager.AddGroup", {"Uid": "foo", "Uidlist": ["bar"]})
