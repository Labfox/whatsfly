from typing import List, Self

STAGE_BETA = 0
STAGE_PROD = 1

class Feature():
    '''
    A class representing a feature.
    You can directly assign the class's variables
    '''
    def __init__(self):
        self.name = None
        self.version = 0
        self.stage = STAGE_PROD
        self.function_name = None


    def compile_feature(self) -> dict:
        '''
        A function to translate a feature to a dict
        :return: A dict containing the translated feature
        '''
        final_string = {
            "name": self.name,
            "version": self.version,
            "stage": "beta" if self.stage == STAGE_BETA else "prod",
            "function_name": self.function_name
        }
        return final_string

    def decompile_feature(self, feature: dict) -> Self:
        '''
        A function to translate a dict into a feature
        :param feature: The feature dict
        '''
        self.name = feature["name"]
        self.version = feature["version"]
        self.stage = STAGE_BETA if feature["stage"] == "beta" else STAGE_PROD
        self.function_name = feature["function_name"]

    def check_feature(self, ofeat: Self) -> bool:
        '''
        Check the compatibilty with another feature (the other feature is considered as a low requirement.)
        :param ofeat: The other feature
        :return: Is the feature compatible ?
        '''
        if ofeat.name != self.name:
            return False
        if ofeat.version > self.version:
            return False
        if ofeat.stage != self.stage:
            return False
        if ofeat.function_name != self.function_name:
            return False
        return True

class Version():
    '''
    A class to represent a binary's version.

    '''
    def __init__(self):
        self._features: List[Feature] = {}
        self._version_type: int = 0



    def add_feature(self, feature: Feature):
        '''
        Add a feature to the version
        :param feature: The feature to add
        '''
        self._features[feature.name] = feature

    def bulk_compile_features(self, feats: dict[Feature]) -> List[dict]:
        '''
        Bulks translates features
        :param feats: The features to translate
        :return: The translated features
        '''
        nfeats = []
        for feat in feats:
            nfeats.append(feats[feat].compile_feature())

        return nfeats

    def bulk_decompile_features(self, feats: List[dict]) -> dict[Feature]:
        '''
        Bulks deserialize features
        :param feats: The features to deserialize
        :return: The deserialized features
        '''
        nfeats = {}
        for feat in feats:
            nfeat = Feature()
            nfeat.decompile_feature(feat)

            nfeats[nfeat.name] = nfeat

        return nfeats

    def compile_version(self) -> dict:
        '''
        Serializes a version
        :return: The serialized version
        '''
        final_data = {
            "features": self.bulk_compile_features(self._features),
            "version_type": self._version_type
        }
        return final_data

    def decompile_version(self, version: dict):
        '''
        Deserializes a version
        :param version: The serialized version
        '''
        if version["version_type"] == 0:
            self._features = self.bulk_decompile_features(version["features"])

    def check_feature(self, feature: Feature) -> bool:
        '''
        Checks if the version contains a feature
        :param: The feature to verify
        :return: Is the feature present
        '''
        if feature.name in self._features:
            return self._features[feature.name].check_feature(feature)
        return False

def return_version():
    pass