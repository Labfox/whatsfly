import unittest
from whatsfly.featureflags import Version, Feature, STAGE_BETA, STAGE_PROD

class FeaturesTests(unittest.TestCase):
    def test_feature_compilation(self):
        feat = Feature()
        feat.name = "Test"
        feat.function_name = "Test"
        feat.stage = STAGE_PROD
        feat.version = 1

        compiled = feat.compile_feature()
        decompiled = Feature()
        decompiled.decompile_feature(compiled)

        self.assertEqual(feat.compile_feature(), decompiled.compile_feature())  # add assertion here
        self.assertEqual(feat.name, decompiled.name)
        self.assertEqual(feat.function_name, decompiled.function_name)
        self.assertEqual(feat.stage, decompiled.stage)
        self.assertEqual(feat.version, decompiled.version)

    def test_feature_comparaison(self):
        feat = Feature()
        feat.name = "Test"
        feat.function_name = "Test"
        feat.stage = STAGE_PROD
        feat.version = 1

        feat2 = Feature()
        feat2.decompile_feature(feat.compile_feature())
        feat2.version = 4

        self.assertTrue(feat2.check_feature(feat))
        self.assertFalse(feat.check_feature(feat2))

class VersionTests(unittest.TestCase):
    def create_test_versions(self):
        feat = Feature()
        feat.name = "Test"
        feat.function_name = "Test"
        feat.stage = STAGE_PROD
        feat.version = 1

        feat2 = Feature()
        feat2.decompile_feature(feat.compile_feature())
        feat2.version = 4

        feat3 = Feature()
        feat3.name = "Invalid"

        return feat, feat2, feat3

    def test_version_compilation(self):
        feat, feat2, feat3 = self.create_test_versions()

        ver = Version()
        ver.add_feature(feat2)
        ver.add_feature(feat3)

        ver2 = Version()
        ver2.decompile_version(ver.compile_version())

        self.assertEqual(ver._version_type, ver2._version_type)
        self.assertEqual(len(ver._features), len(ver2._features))

    def test_version_check(self):
        feat, feat2, feat3 = self.create_test_versions()

        ver = Version()
        ver.add_feature(feat2)
        ver.add_feature(feat3)





        self.assertTrue(ver.check_feature(feat))
        self.assertTrue(ver.check_feature(feat3))

        del ver

        ver1 = Version()
        ver1.add_feature(feat)

        self.assertFalse(ver1.check_feature(feat2))
        self.assertFalse(ver1.check_feature(feat3))
