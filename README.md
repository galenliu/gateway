本项目参考了[mozilla WebThins gateway]( https://iot.mozilla.org/gateway/#carousel-section2 )
项目的相关实现。 因为在学习mozilla WebThings平台时，由于众所周知的原因，Nodejs依赖包无法正常工作，于是这成了原始动机！

### API说明：

- Gateway API：


- Plugin API如：T Get: /plugin


- Adapter API：

{GET /new_things HTTP/1.1 1 1 map[Accept-Encoding:[gzip, deflate, br]
Accept-Language:[zh-CN,zh;q=0.9,zh-HK;q=0.8] Cache-Control:[no-cache] Connection:[Upgrade]
Cookie:[_ga=GA1.1.80962085.1605737746; ___rl__test__cookies=1610012956392] Origin:[http://localhost:3000] 
Pragma:[no-cache] Sec-Websocket-Extensions:[permessage-deflate; client_max_window_bits]
Sec-Websocket-Key:[ojWvG3ByY6oqUxLWqnTkzw==] Sec-Websocket-Version:[13] Upgrade:[websocket]
User-Agent:[Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36]]
{} <nil> 0 [] false localhost:9090 map[] map[] <nil> map[] [::1]:53000 /new_things <nil> <nil> <nil>
0xc0004320c0}

{GET /things/ HTTP/1.1 1 1 map[Accept-Encoding:[gzip, deflate, br] Accept-Language:[zh-CN,zh;q=0.9,zh-HK;q=0.8]
Cache-Control:[no-cache] Connection:[Upgrade] Cookie:[_ga=GA1.1.80962085.1605737746; ___rl__test__cookies=1610012956392]
Origin:[http://localhost:3000] Pragma:[no-cache] Sec-Websocket-Extensions:[permessage-deflate; client_max_window_bits]
Sec-Websocket-Key:[gBkpB2aerMWj9PfXIowg8A==] Sec-Websocket-Version:[13] Upgrade:[websocket]
User-Agent:[Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36]]
{} <nil> 0 [] false localhost:9090 map[] map[] <nil> map[] [::1]:53482 /things/ <nil> <nil> <nil> 0xc000433240}