### How to use
- ```mkdir /root/redi301 && cd /root/redi301 && wget https://raw.githubusercontent.com/haha357/redi301/master/bin/redi301 && chmod +x ./redi301```
- ```./redi301```

### Instructions
  This application is aimed at redirecting urls which have been blocked by GFW to these which have not been blocked, this is a real 301 redirect 
  application, because it responds 301 status code in the response headers. 

  Use```./redi301 -h``` to view more instructions. 

### Params
- -a Listen address，just like ```ip:port```, exp : ```0.0.0.0:80```, the default value is ```0.0.0.0:80```.
- -t Prefix of the target redirect url, just like：```http://www.whitehouse.gov```,both http and https supported, please do not end with ```/```.

### TODO
- [x] Http redirect to Http. (exp: http://blocked.domain -> http://unblocked.domain )
- [x] Http redirect to Https. (exp: http://blocked.domain -> https://unblocked.domain )
- [ ] Https redirect to Http. (exp: https://blocked.domain -> http://unblocked.domain )
- [ ] Https redirect to Https. (exp: https://blocked.domain -> https://unblocked.domain )

### Donation

- If you like this application, you can donate to the author. Thank you so much.
- USDT(TRC20) wallet: TB8meT4Pm9KFXRJ8SNCfxx4yBGPbk3Ekip

  <div style="text-align: center; width: 500px; border: green solid 1px;"><img src="https://img.mdev.eu.org/file/5bda398b80a9ce195b72c.png"></div>

### itdog.cn test

<img src="https://img.mdev.eu.org/file/3264671ad7801e86649e7.png">


### BenchMark
