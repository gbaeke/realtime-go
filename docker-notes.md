Dockerfile uses following command in first stage:

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

Note that it this creates a passwd & group file in the /user folder with just the nobody user and group for later copying to the scratch image.

passwd format is: user:passwd:userid:groupid:userid-info:homedir:cmd-shell

Note: password is x to indicate the password is stored in shadow password file

Note: be aware that you will NOT be able to bind to low ports (under 1024) when running as nobody; if you run as non root and try to bind low port then:

2018/12/18 17:38:44 PONG <nil>
2018/12/18 17:38:44 Serving on localhost:8888...
2018/12/18 17:38:44 listen tcp :88: bind: permission denied