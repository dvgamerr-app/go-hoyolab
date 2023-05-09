# Hoyolab Daily Checkin
[![Go](https://github.com/dvgamerr/go-hoyolab/actions/workflows/build.yml/badge.svg)](https://github.com/dvgamerr/go-hoyolab/actions/workflows/build.yml)
[![CodeQL](https://github.com/dvgamerr/go-hoyolab/actions/workflows/codeql.yml/badge.svg)](https://github.com/dvgamerr/go-hoyolab/actions/workflows/codeql.yml)
[![Dependency Review](https://github.com/dvgamerr/go-hoyolab/actions/workflows/review.yml/badge.svg)](https://github.com/dvgamerr/go-hoyolab/actions/workflows/review.yml)

Genshin Impact's Hoyolab Daily Check-in Bot is here!. You only need to run it once, then it will continue to run forever.
![example.png](./docs/example.png)

### In Progress
- [ ] Honkai StarRail cliams reward support.
- [ ] Honkai Impact 3 cliams reward support.

## How to use
1. Open chrome browser open [https://www.hoyolab.com/home](https://www.hoyolab.com/home)
2. Login user genshin account for daily cliams.
3. run `hoyolab.exe`.
4. If found Error please craete issues in [https://github.com/dvgamerr/go-hoyolab/issues](https://github.com/dvgamerr/go-hoyolab/issues)

## Solution for Automatic
- use `Task Scheduler` and `Create Basic Task`
- select `When I log on` If you turn on Computer every day, or Select `Daily` If your run cumputer 24/7.

![image](https://user-images.githubusercontent.com/10203425/236996927-cb76c5be-09be-409c-8cb2-743bb0204d1a.png)

- after click Next select `Start a program`, Next and browse `hoyolab.exe`

## Todos
- [ ] update new version from github automatic
- [ ] install schedule task with windows-os automatic.
- [ ] command line `hoyolab` support all os.
- [ ] support session with all browser.
- [ ] docker container support

## Done
- [x] cliams checkin daily with chrome session.
- [x] all profile chrome session.

## Prerequisites
- Windows OS
- Have login to mihoyo's website at any browser (A login for a year is enough)

## License
MIT © Touno.io 2023
