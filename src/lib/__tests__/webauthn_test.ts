import {
    PublicKeyCredentialCreationOptions as CCO,
    PublicKeyCredentialParameters,
    PublicKeyCredentialUserEntity,
    PublicKeyCredentialRpEntity,
    AuthenticatorSelectionCriteria,
} from "../../protos/webauthn_pb";

import { creationOptionsFromProto } from "../webauthn";

import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

test("can do a basic conversion", () => {
    const input = new CCO();

    const challenge = Uint8Array.from("challenge", (c) => c.charCodeAt(0));
    input.setChallenge(challenge);

    const cred = new PublicKeyCredentialParameters();
    cred.setAlg(-1);
    input.addPubKeyCredParams(cred);

    const timeout = new google_protobuf_duration_pb.Duration();
    timeout.setSeconds(60);
    input.setTimeout(timeout);

    const user = new PublicKeyCredentialUserEntity();
    user.setName("user");
    user.setDisplayName("a user");
    const userID = Uint8Array.of(1, 2, 3, 4);
    user.setId(userID);
    input.setUser(user);

    const rp = new PublicKeyCredentialRpEntity();
    rp.setName("foo");
    rp.setId("foo id");
    input.setRp(rp);

    const auths = new AuthenticatorSelectionCriteria();
    auths.setUserVerification(
        AuthenticatorSelectionCriteria.UserVerificationRequirement.DISCOURAGED
    );
    auths.setAuthenticatorAttachment(
        AuthenticatorSelectionCriteria.AuthenticatorAttachment.CROSS_PLATFORM
    );
    input.setAuthenticatorSelection(auths);

    const want: PublicKeyCredentialCreationOptions = {
        challenge: challenge,
        excludeCredentials: [],
        pubKeyCredParams: [
            {
                alg: -1,
                type: "public-key",
            },
        ],
        timeout: 60000,
        rp: {
            id: "foo id",
            name: "foo",
        },
        user: {
            name: "user",
            displayName: "a user",
            id: userID,
        },
        authenticatorSelection: {
            authenticatorAttachment: "cross-platform",
            userVerification: "discouraged",
            requireResidentKey: false,
        },
    };

    const got = creationOptionsFromProto(input);
    expect(got).toStrictEqual(want);
});
