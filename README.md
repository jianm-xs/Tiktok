## Tiktok（青训营——可能不太队）

[![All Contributors](https://img.shields.io/badge/all_contributors-5-orange.svg?style=flat-square)](#contributors-)![GitHub closed pull requests](https://img.shields.io/github/issues-pr-closed/jianm-xs/Tiktok?color=brightgreen)![GitHub closed issues](https://img.shields.io/github/issues-closed/jianm-xs/Tiktok?color=brightgreen)![GitHub forks](https://img.shields.io/github/forks/jianm-xs/Tiktok?color=cyan)![GitHub watchers](https://img.shields.io/github/watchers/jianm-xs/Tiktok?color=cyan)![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/jianm-xs/Tiktok/main?color=blue&filename=Project%2Fgo.mod)

第三届字节跳动青训营——后端专场**可能不太队**极简抖音项目

### 上手指南

以下指南将帮助你在本地机器上安装和运行项目，进行开发和测试。

##### 环境要求

- [FFmpeg](https://ffmpeg.org/) 
- [Go1.18](https://golang.google.cn/)
- [MySQL](https://www.mysql.com/)
- [Git](https://git-scm.com/)
- [Redis](https://redis.io/)

##### 安装步骤

1. clone 本项目：`git clone https://github.com/jianm-xs/Tiktok.git `
2. 建立数据库，运行 `数据库文件/Tiktok.sql` 文件
3. 进入项目：`cd Tiktok/Project`
4. 修改项目配置：`vim /config/config.go`
   - 修改 `MysqlCfg` 为自己的数据库信息
   - 修改 `RedisCfg` 为自己的 Redis 信息
   - 修改 `ServerHost` 为自己的服务器地址（用于访问上传的文件）
   - 修改 `ServerPort` 为自己的服务端口
5. 编译项目：`go build`
6. 运行项目：`./Project`

##### App 端

使用青训营提供的 Android 程序，说明地址：[极简抖音 App](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)

### 模块功能

|                    接口                     |               说明               |
| :-----------------------------------------: | :------------------------------: |
|         视频流接口（/douyin/feed/）         |   返回视频列表，单次最多 30 个   |
|     用户注册（/douyin/user/register/）      |   注册新用户，注册成功自动登录   |
|       用户登录（/douyin/user/login/）       | 用户登录接口，返回用户鉴权 token |
|     视频投稿（/douyin/publish/action/）     |           用户发布视频           |
|          用户信息（/douyin/user/）          |        查看用户的所有信息        |
|      发布列表（/douyin/publish/list/）      |      查看用户发布的所有视频      |
|    点赞操作（/douyin/favorite/action/）     |          对视频进行点赞          |
|     点赞列表（/douyin/favorite/list/）      |      查看用户点赞的所有视频      |
|     评论操作（/douyin/comment/action/）     |        用户对视频发布评论        |
|      评论列表（/douyin/comment/list/）      |        查看视频的所有评论        |
|    关注操作（/douyin/relation/action/）     |           关注某一用户           |
|  关注列表（/douyin/relation/follow/list/）  |        查看某用户关注的人        |
| 粉丝列表（/douyin/relation/follower/list/） |         查看某用户的粉丝         |

详情可查看 [API 文档](https://www.apifox.cn/apidoc/shared-dbc54832-2446-428e-88a0-05f2a7e42250)

### 如果要参与贡献，请仔细阅读 `贡献文档.md`。

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<table>
  <tr>
    <td align="center"><a href="https://github.com/jianm-xs"><img src="https://avatars.githubusercontent.com/u/69761545?v=4?s=100" width="100px;" alt=""/><br /><sub><b>jianm-xs</b></sub></a><br /><a href="https://github.com/jianm-xs/Tiktok/commits?author=jianm-xs" title="Documentation">📖</a> <a href="#tutorial-jianm-xs" title="Tutorials">✅</a> <a href="#business-jianm-xs" title="Business development">💼</a> <a href="https://github.com/jianm-xs/Tiktok/commits?author=jianm-xs" title="Code">💻</a> <a href="#projectManagement-jianm-xs" title="Project Management">📆</a> <a href="https://github.com/jianm-xs/Tiktok/issues?q=author%3Ajianm-xs" title="Bug reports">🐛</a> <a href="#question-jianm-xs" title="Answering Questions">💬</a></td>
    <td align="center"><a href="https://github.com/LuWiHan"><img src="https://avatars.githubusercontent.com/u/96118540?v=4?s=100" width="100px;" alt=""/><br /><sub><b>LuWiHan</b></sub></a><br /><a href="https://github.com/jianm-xs/Tiktok/commits?author=LuWiHan" title="Documentation">📖</a> <a href="#design-LuWiHan" title="Design">🎨</a> <a href="https://github.com/jianm-xs/Tiktok/commits?author=LuWiHan" title="Code">💻</a></td>
    <td align="center"><a href="https://gitee.com/wrz0318"><img src="https://avatars.githubusercontent.com/u/74159645?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Nicknamezz00</b></sub></a><br /><a href="https://github.com/jianm-xs/Tiktok/commits?author=Nicknamezz00" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/KiceAmber"><img src="https://avatars.githubusercontent.com/u/90232365?v=4?s=100" width="100px;" alt=""/><br /><sub><b>kice</b></sub></a><br /><a href="https://github.com/jianm-xs/Tiktok/commits?author=KiceAmber" title="Code">💻</a> <a href="#question-KiceAmber" title="Answering Questions">💬</a></td>
    <td align="center"><a href="https://github.com/bingguoq"><img src="https://avatars.githubusercontent.com/u/103885711?v=4?s=100" width="100px;" alt=""/><br /><sub><b>bingguoq</b></sub></a><br /><a href="https://github.com/jianm-xs/Tiktok/commits?author=bingguoq" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!