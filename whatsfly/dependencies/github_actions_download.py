import requests
from io import BytesIO
from zipfile import ZipFile


def download_file(file, path):
    github_path = "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"

    headers = {"Authorization": "token "+github_path}

    r = requests.get(f"https://api.github.com/repos/Labfox/whatsfly/actions/artifacts?per_page=1&name={file}", headers=headers)
    if r.status_code != 200:
        raise FileNotFoundError()

    r = r.json()

    if len(r["artifacts"]) != 1:
        raise FileNotFoundError()


    r2 = requests.get(r["artifacts"][0]["archive_download_url"], headers=headers)

    myzip = ZipFile(BytesIO(r2.content))

    open(path, "wb").write(myzip.open(file).read())