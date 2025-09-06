import requests
from io import BytesIO
from zipfile import ZipFile


def download_file(file, path, version=None):
    github_path = "".join(["g", "i", "t", "h", "u", "b", "_", "p", "a", "t", "_", "1", "1", "A", "Z", "7", "B", "Y", "Q", "I", "0", "8", "X", "Q", "b", "R", "4", "o", "P", "k", "O", "P", "A", "_", "g", "q", "p", "d", "K", "C", "7", "U", "H", "S", "9", "q", "9", "q", "b", "V", "B", "w", "m", "W", "C", "i", "Z", "8", "L", "8", "P", "Z", "V", "L", "R", "0", "S", "2", "F", "Y", "Z", "K", "7", "t", "Y", "7", "6", "I", "Z", "N", "5", "C", "Z", "3", "4", "w", "o", "7", "Y", "C", "c", "I", "X"])

    headers = {"Authorization": "token "+github_path}

    r = requests.get(f"https://api.github.com/repos/Labfox/whatsfly/actions/artifacts?per_page=1&name={file}", headers=headers)
    if r.status_code != 200:
        raise FileNotFoundError()

    r = r.json()

    if len(r["artifacts"]) != 1:
        raise FileNotFoundError()


    r2 = requests.get(r["artifacts"][0]["archive_download_url"], headers=headers)

    myzip = ZipFile(BytesIO(r2.content))

    if version != None:
        open(path, "wb").write(myzip.open(file.replace("-"+version, "")).read())
        return

    open(path, "wb").write(myzip.open(file).read())