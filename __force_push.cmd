@echo off
chcp 65001
set GIT_COMMITTER_DATE=2025-12-20T12:31:58
git add .
git -c user.name="想太多" -c user.email="huangyuanlei@hotmail.com" commit --amend --date=2025-12-20T12:31:58 --no-edit
git push --force
