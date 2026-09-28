@echo off

cd /d %~dp0

@REM git init
@REM git config user.email "bingk3069@gmail.com"
@REM git config user.name "dx838"

@REM git add *
@REM git commit -m "first commit"
@REM git branch -M main
@REM git remote add origin https://github.com/dx838/BiliCacheMerge.git
@REM git push -u origin main





@REM git remote add origin https://github.com/dx838/BiliCacheMerge.git
@REM git remote set-url origin https://github.com/dx838/BiliCacheMerge.git
@REM git branch -M main


git add *
git commit -m "Test 0 %date% %time%"
git push -u origin main



rem java keytool.exe -genkeypair -v -storetype JKS -keyalg RSA -keysize 2048 -validity 10000   -keystore bbmusic-keystore.jks   -alias bbmusic   -keypass passwd   -storepass passwd   -dname "CN=GD, OU=CN, O=None, L=MM, ST=NN, C=CN"

rem base64 -i  bbmusic-keystore.jks


pause

rem 删除损坏的索引
rem rm .git/index
rem git reset


rem 强制用本地覆盖远程（谨慎用）
rem  git pull origin main --allow-unrelated-histories
rem  git push origin main -f

rem Git 强制用远程覆盖本地（彻底同步）
rem rem git fetch --all && git reset --hard origin/$(git branch --show-current)
rem git fetch --all
rem git reset --hard origin/main


rem  创建无历史的孤儿分支，再强制覆盖主分支。
rem # 2. 创建孤儿分支（无任何历史）
rem git checkout --orphan new-empty

rem # 3. 清空所有文件
rem git rm -rf .
rem echo "# 清空后初始化" > README.md
rem git add .
rem git commit -m "None"

rem # 4. 强制推送到远程主分支
rem git push -f origin new-empty:main

rem # 5. 切换回主分支并清理
rem rem git checkout main
rem rem git pull origin main
rem git branch -M main
rem git branch -D new-empty

