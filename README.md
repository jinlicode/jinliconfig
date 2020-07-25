# jinliconfig

锦鲤容器管理部署工具是致力于解决中小站长部署环境问题

## 优势

1. 使用docker容器化技术进行文件系统隔离
2. 除了前端openresty服务器需要暴露80和443端口，其他所有服务都运行在容器子网
3. 支持一键脚本部署，仅需要一个命令即可部署完成
4. 抛弃了传统管理面板占用服务器资源，不运行不需要占用服务器任何资源
5. 支持服务器远程终端Tui图形化管理，方便友好快捷，仅需要记住一个命令 **`jinliconfig`** 即可

## 产品功能规划

此部分为协作开发者使用，普通用户作为了解既可，这部分也详细讲解了锦鲤部署的内部运行机制，欢迎大家监督考察，有能力可以阅读代码，和贡献代码， **因为开源，所以伟大**

### v1.0版本

此版本为初始发行版本支持的功能规划版本

#### 初始化部署

```flow
st=>start: Start:>https://www.markdown-syntax.com
io=>inputoutput: verification
op=>operation: Your Operation
cond=>condition: Yes or No?
sub=>subroutine: Your Subroutine
e=>end
st->io->op->cond
cond(yes)->e
cond(no)->sub->io
```

#### 管理部署

