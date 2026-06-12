import smtplib
from email.mime.text import MIMEText

def SendMessageSMTP(msg :str):
    msg = MIMEText(msg)
    msg["Subject"] = "Test"
    msg["From"] = "renta64@yandex.ru"
    msg["To"] = "zenin.nickita2013@yandex.ru"

    server = smtplib.SMTP_SSL("smtp.yandex.ru", 465)
    server.login("renta64@yandex.ru", "ltfldquzanehugwy")
    server.send_message(msg)
    server.quit()