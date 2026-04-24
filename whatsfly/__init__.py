import os
import sys
import logging
import time

from .dependencies.builder import ensureUsableBinaries

try:
    if not ("WHATSFLY_NO_UPDATES" in os.environ):
        root_dir = os.path.abspath(os.path.dirname(__file__))
        update_file = os.path.join(root_dir, "last_binary_update.txt")
        if os.path.exists(update_file):
            last_update = open(update_file, "r").read()
            if time.time() - int(last_update) > 60 * 60 * 24 * 31:
                from sys import platform

                if platform == "darwin":
                    file_ext = "latest.dylib"
                elif platform in ("win32", "cygwin"):
                    file_ext = "latest.dll"
                else:
                    file_ext = "latest.so"
                os.remove(update_file)
                os.remove(f"{root_dir}/dependencies/{file_ext}")
        else:
            f = open(update_file, "w")
            f.write(str(int(time.time())))
            f.close()
except Exception as e:
    print("Unable to ensure updates timeframe")
    print(e)

try:
    from .whatsapp import WhatsApp
except OSError:
    ensureUsableBinaries()
    from .whatsapp import WhatsApp

LOGGER = logging.getLogger()
logging.basicConfig(level=logging.INFO)
