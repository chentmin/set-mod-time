# 处理文件最后更改时间

根据文件sha1 checksum, 找出没有变化的文件, 把文件的最后更新时间改为之前的值, 确保文件内容没有改动时, 最后更新时间也相同

aws s3 sync 或者 rsync 等同步工具, 根据文件大小和文件最后更新时间来判断文件是否有变化. 每次git clone下来的文件的最后更新时间都是当前时间. 所以需要本工具

# 用法

可结合travis一起使用, 用travis来缓存数据文件

`.travis.yml`中, 缓存数据文件

    cache:
      directories:
      - $HOME/.config

下载文件并运行

    wget https://raw.githubusercontent.com/szmolin/dist/master/set-mod-time -O $HOME/set-mod-time
    chmod +x $HOME/set-mod-time
    $HOME/set-mod-time -config=$HOME/.config/config.txt -path=.deploy