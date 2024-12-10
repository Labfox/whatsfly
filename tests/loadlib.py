import os

os.system('alias go="echo e"')

from whatsfly import WhatsApp


client = WhatsApp(on_event=print)
