import {
    PublicKeyCredentialCreationOptions as CCO,
    PublicKeyCredentialRequestOptions as CRO,
    PublicKeyCredentialParameters,
    PublicKeyCredentialUserEntity,
    PublicKeyCredentialRpEntity,
    AuthenticatorSelectionCriteria,
    PublicKeyCredential as C,
    PublicKeyCredentialDescriptor,
    AuthenticatorResponse as AR,
    AuthenticatorAssertionResponse as AAsR,
    AuthenticatorAttestationResponse as AAtR,
} from "../../protos/webauthn_pb";

import { creationOptionsFromProto, credentialFromJS, requestOptionsFromProto } from "../webauthn";

import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";

test("can convert a creation request in proto form to javascript form", () => {
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

    const excludedCredential = new PublicKeyCredentialDescriptor();
    excludedCredential.setId(Uint8Array.of(1, 2, 3, 4));
    excludedCredential.setType("public-key");
    excludedCredential.addTransports(PublicKeyCredentialDescriptor.AuthenticatorTransport.BLE);
    excludedCredential.addTransports(PublicKeyCredentialDescriptor.AuthenticatorTransport.INTERNAL);
    input.addExcludeCredentials(excludedCredential);

    const want: PublicKeyCredentialCreationOptions = {
        challenge: challenge,
        excludeCredentials: [
            {
                id: Uint8Array.of(1, 2, 3, 4),
                type: "public-key",
                transports: ["ble", "internal"],
            },
        ],
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

test("can convert an attestation response to a proto", () => {
    const input: PublicKeyCredential = {
        id: "abc",
        type: "public-key",
        rawId: Uint8Array.from("abc", (c) => c.charCodeAt(0)),
        response: {
            clientDataJSON: Uint8Array.from("{}", (c) => c.charCodeAt(0)),
            attestationObject: Uint8Array.from("foo", (c) => c.charCodeAt(0)),
        } as AuthenticatorAttestationResponse,
        getClientExtensionResults: () => {
            return {};
        },
    };
    const want = new C();
    want.setId("abc");
    want.setType("public-key");
    const r = new AR();
    r.setClientDataJson(Uint8Array.from("{}", (c) => c.charCodeAt(0)));
    const atr = new AAtR();
    atr.setAttestationObject(Uint8Array.from("foo", (c) => c.charCodeAt(0)));
    r.setAttestationResponse(atr);
    want.setResponse(r);

    const got = credentialFromJS(input);
    expect(got.toObject()).toStrictEqual(want.toObject());
});

test("can convert an assertion response to a proto", () => {
    const input: PublicKeyCredential = {
        id: "abc",
        type: "public-key",
        rawId: Uint8Array.from("abc", (c) => c.charCodeAt(0)),
        response: {
            clientDataJSON: Uint8Array.from("{}", (c) => c.charCodeAt(0)),
            authenticatorData: Uint8Array.from("auth", (c) => c.charCodeAt(0)),
            signature: Uint8Array.from("sig", (c) => c.charCodeAt(0)),
            userHandle: Uint8Array.from("uh", (c) => c.charCodeAt(0)),
        } as AuthenticatorAssertionResponse,
        getClientExtensionResults: () => {
            return {};
        },
    };
    const want = new C();
    want.setId("abc");
    want.setType("public-key");
    const r = new AR();
    r.setClientDataJson(Uint8Array.from("{}", (c) => c.charCodeAt(0)));
    const asr = new AAsR();
    asr.setAuthenticatorData(btoa("auth"));
    asr.setSignature(btoa("sig"));
    asr.setUserHandle(btoa("uh"));
    r.setAssertionResponse(asr);
    want.setResponse(r);

    const got = credentialFromJS(input);
    expect(got.toObject()).toStrictEqual(want.toObject());
});

test("can convert a credential request proto to an object", () => {
    const input = new CRO();
    input.setChallenge(btoa("foo"));

    const timeout = new google_protobuf_duration_pb.Duration();
    timeout.setSeconds(60);
    input.setTimeout(timeout);

    const credential = new PublicKeyCredentialDescriptor();
    credential.setId(Uint8Array.of(1, 2, 3, 4));
    credential.setType("public-key");
    credential.addTransports(PublicKeyCredentialDescriptor.AuthenticatorTransport.BLE);
    credential.addTransports(PublicKeyCredentialDescriptor.AuthenticatorTransport.INTERNAL);
    input.addAllowedCredentials(credential);

    const want: PublicKeyCredentialRequestOptions = {
        challenge: Uint8Array.from("foo", (c) => c.charCodeAt(0)),
        timeout: 60000,
        allowCredentials: [
            {
                id: Uint8Array.of(1, 2, 3, 4),
                type: "public-key",
                transports: ["ble", "internal"],
            },
        ],
    };
    const got = requestOptionsFromProto(input);
    expect(got).toStrictEqual(want);
});
