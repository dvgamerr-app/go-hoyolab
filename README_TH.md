<div align="center">
  <h1>แอป Check-in อัตโนมัติสำหรับ Hoyoverse Game</h1>
  <p>
    <a href="https://github.com/dvgamerr/go-hoyolab/actions/workflows/build.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/dvgamerr/go-hoyolab/build.yml?label=Build&amp;style=flat-square" alt="GitHub Build Action Status">
    </a>
    <img src="https://img.shields.io/github/actions/workflow/status/dvgamerr/go-hoyolab/review.yml?label=Dependency&amp;style=flat-square" alt="GitHub Build Action Status">
    <a href="https://goreportcard.com/report/dvgamerr/go-hoyolab">
      <img src="https://goreportcard.com/badge/dvgamerr/go-hoyolab?style=flat-square">
    </a>
    <br>
    <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square">
    <a href="LICENSE.md">
      <img src="https://img.shields.io/github/license/dvgamerr/go-hoyolab?style=flat-square" alt="LICENSE">
    </a>
    <a href="https://github.com/dvgamerr/go-hoyolab/releases/latest">
      <img src="https://img.shields.io/github/release-date/dvgamerr/go-hoyolab?style=flat-square">
    </a>
  </p>
  <p><a href="./README.md">English</a> | <a href="./README_TH.md">ภาษาไทย</a></p>
  <p>Genshin Impact, Honkai StarRail, Honkai Impact 3. You only need to run it once, then it will continue to run forever.</p>
</div>



![example.png](./docs/example-logs.png)

![checkin.png](./docs/checkin.png)

## Features
- [x] ปล่อย version ใหม่ผ่าน github actions.
- [x] เริ่มต้นใช้งานได้แค่ chrome เท่านั้น.
- [x] สามารถใช้กับ chrome หลายโปรไฟล์ เพื่อใช้กับหลาย login ID ได้.
- [X] สามารถรัน `hoyolab` ได้ทุก os.
- [x] แจ้งเตือนผ่าน Line Notify.

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/0cbdb857-f866-4813-8420-03c2ce73688e)

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/133f8fcd-d301-471f-92a7-6e88874ff851)

- [x] แจ้งเตือนผ่าน discord webhook.

![checkin.png](./docs/example-notify.png)


  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/1c75dc54-e787-4831-94a0-047f1aef7e1a)
  
- [x] รับรองเกม Genshin Impact.
- [x] รับรองเกม Honkai StarRail.
- [x] รับรองเกม Honkai Impact 3.

### ที่จะทำต่อ
- [ ] รองรับ Zenless Zone Zero.
- [ ] docker container support.
- [ ] support session with all browser.
- [ ] install schedule task with windows-os automatic.

---

## วิธีใช้
1. เปิด chrome ขึ้นมาแล้วเข้าเว็บ [https://www.hoyolab.com/home](https://www.hoyolab.com/home)
2. เข้าใช้งาน ด้วย account ของ hoyoverse ที่ต้องการจะรับรางวัน.
3. เปิดโปรแกรม `hoyolab.exe` ขึ้นมา เสร็จ.
4. ถ้าพบ error แจ้งได้ที่  [https://github.com/dvgamerr/go-hoyolab/issues](https://github.com/dvgamerr/go-hoyolab/issues)
5. ถ้าต้องการ ใช้แจ่้งเตือนให้เข้่าไปที่ [LINE-Notify](https://notify-bot.line.me/my/)
6. เข้าใข้งานแล้ว กด `Generate token` จะได้ token มาแล้ว copy มาวางไว้ใน `hoyolab.yaml` ที่ `notify.token` แทนที่ `XXXXXX` ตามรูป:
 
 ![image](./docs/example-token.png)
  
7.ถ้าเกมไหน ไม่ได้เล่นให้ลบ หรือ `#` เพื่อปิดไปได้เลยตามรูป

  ![image](https://github.com/dvgamerr/go-hoyolab/assets/10203425/7ab44d88-31cf-4919-ab5a-e7c4da5beedf)



### Windows
### สำหรับเครื่องที่เปิดเครื่อง 24 ชม ให้สร้าง Task Schedule ตามนี้ (แนะนำให้ใช้ profile อื่นเพื่อแทน profile หลัก)
- เปิด `Task Scheduler` คลิก `Create Basic Task`
- เลือก `Daily` แล้วกด next เลือกเวลาเป็น `5 am`.
- หลังจากกด Next ให้เลือก `Start a program`, กด next และเลือก browse ไปหาไฟล์ `hoyolab.exe` ในเครื่อง

### สำหรับเครื่องที่เปิด-ปิดทุกๆ วัน
- สร้าง shortcut จากไฟล์ `hoyolab.exe`
- แล้ว copy ไปวางไว้ที่ `C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup`
- ทุกวันที่เปิดเครื่องมันจะ checkin ให้ ตลอด 

### MaOcOS & Linux
- ใช้ `crontab` ได้เลย ง่ายกว่าเยอะ.

## Prerequisites
- เข้าสู่ระบบเว็บไซต์ของ mihoyo ด้วยเบราว์เซอร์ Chrome

## หากคุณต้องการความช่วยเหลือ คุณสามารถเข้าร่วม Discord ได้
#### (อย่าเที่ยวพิมพ์ @everyone แล้ว โวยวายเป็นเด็กๆ อธิบายดีๆ เด๋วมาตอบ)

[![Join Us?](https://discordapp.com/api/guilds/475720106471849996/widget.png?style=banner2)](https://discord.gg/QDccF497Mw)
