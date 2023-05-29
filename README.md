<div align="center">
  <h1>Hoyolab a Daily Check-in for Hoyoverse Game</h1>
  <p>
    <a href="https://github.com/dvgamerr/go-hoyolab/actions/workflows/build.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/dvgamerr/go-hoyolab/build.yml?label=Build&amp;style=flat-square" alt="GitHub Build Action Status">
    </a>
    <a href="https://github.com/dvgamerr/go-hoyolab/actions/workflows/codeql.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/dvgamerr/go-hoyolab/codeql.yml?label=CodeQL&amp;style=flat-square" alt="GitHub Build Action Status">
    </a>
    <a href="https://github.com/dvgamerr/go-hoyolab/actions/workflows/review.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/dvgamerr/go-hoyolab/review.yml?label=Dependency&amp;style=flat-square" alt="GitHub Build Action Status">
    </a>
    <a href="LICENSE.md">
      <img src="https://img.shields.io/github/license/dvgamerr/go-hoyolab?style=flat-square" alt="LICENSE">
    </a>
    <a href="https://github.com/dvgamerr/go-hoyolab/releases/latest">
      <img src="https://img.shields.io/github/release-date/dvgamerr/go-hoyolab?style=flat-square" alt="GitHub Release Date - Published_At">
    </a>
  </p>
  <p>Genshin Impact, Honkai StarRail, Honkai Impact 3. You only need to run it once, then it will continue to run forever.</p>
</div>

![example.png](./docs/example-logs.jpg)

## Features
- [x] Line Notify support after checkin.

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/0cbdb857-f866-4813-8420-03c2ce73688e)

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/133f8fcd-d301-471f-92a7-6e88874ff851)
  
- [x] Claims checkin daily with chrome browser only.
- [x] Multiple chrome profiles for multiple game accounts.

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/1c75dc54-e787-4831-94a0-047f1aef7e1a)
  
- [x] Genshin Impact claims support.
- [x] Honkai StarRail claims support.
- [x] Honkai Impact 3 claims support.

### Todo
- [ ] update new version from github automatic.
- [ ] docker container support.
- [ ] support session with all browser.
- [ ] command line `hoyolab` support all os.
- [ ] install schedule task with windows-os automatic.

---

## How to use
1. Open chrome browser open [https://www.hoyolab.com/home](https://www.hoyolab.com/home)
2. Login user genshin account for daily cliams.
3. run `hoyolab.exe`.
4. If found Error please craete issues in [https://github.com/dvgamerr/go-hoyolab/issues](https://github.com/dvgamerr/go-hoyolab/issues)
5. If Notify message after CheckIn use [LINE-Notify](https://notify-bot.line.me/my/)
6. log in that and `Generate token` in below page after that copy token paste in `hoyolab.yaml` in `notify.token` at `XXXXXX` same image:
 
 ![image](./docs/example-token.png)
  
7. If you don't play some game add `#` in first char in line scope, e.g. i don;t play honkai impact 

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/7ab44d88-31cf-4919-ab5a-e7c4da5beedf)



### Windows
### Solution If Turn On computer 24*7
- use `Task Scheduler` and `Create Basic Task`
- select `Daily` next time at `5 am`.
- after click Next select `Start a program`, Next and browse `hoyolab.exe`

### Solution If starup computer 
- create shortcut from `hoyolab.exe`
- copy shortcut to `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup`

### MaOcOS & Linux
- use `crontab` for automatic run.

## Prerequisites
- Have login to mihoyo's website at chrome browser (A login for a year is enough)

## If you need help You can join Discord.

[![Join Us?](https://discordapp.com/api/guilds/475720106471849996/widget.png?style=banner2)](https://discord.gg/QDccF497Mw)
