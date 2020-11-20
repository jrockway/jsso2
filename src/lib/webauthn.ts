import {
    PublicKeyCredentialCreationOptions as CCO,
    PublicKeyCredentialRequestOptions as CRO,
    PublicKeyCredentialDescriptor as CD,
    AuthenticatorSelectionCriteria as ASC,
    PublicKeyCredential as C,
    AuthenticatorResponse as AR,
    AuthenticatorAssertionResponse as AAsR,
    AuthenticatorAttestationResponse as AAtR,
} from "../protos/webauthn_pb";

function credentialsFromProto(input: CD[]): PublicKeyCredentialDescriptor[] {
    const result = [];
    for (const cred of input) {
        const out = new Object() as PublicKeyCredentialDescriptor;
        out.id = cred.getId_asU8();
        if (cred.getType() == "public-key") {
            out.type = "public-key";
        }
        out.transports = [];
        for (const transport of cred.getTransportsList()) {
            switch (transport) {
                case CD.AuthenticatorTransport.BLE:
                    out.transports.push("ble");
                    break;
                case CD.AuthenticatorTransport.INTERNAL:
                    out.transports.push("internal");
                    break;
                case CD.AuthenticatorTransport.NFC:
                    out.transports.push("nfc");
                    break;
                case CD.AuthenticatorTransport.USB:
                    out.transports.push("usb");
                    break;
            }
        }
        result.push(out);
    }
    return result;
}

export function creationOptionsFromProto(rawOpts: CCO): PublicKeyCredentialCreationOptions {
    const opts = new Object() as PublicKeyCredentialCreationOptions;

    switch (rawOpts.getAttestation()) {
        case CCO.AttestationConveyancePreference.DIRECT:
            opts.attestation = "direct";
            break;
        case CCO.AttestationConveyancePreference.INDIRECT:
            opts.attestation = "indirect";
            break;
    }
    if (rawOpts.hasAuthenticatorSelection()) {
        const auths = rawOpts.getAuthenticatorSelection();
        opts.authenticatorSelection = new Object() as AuthenticatorSelectionCriteria;
        opts.authenticatorSelection.requireResidentKey = auths.getRequireResidentKey();
        switch (auths.getAuthenticatorAttachment()) {
            case ASC.AuthenticatorAttachment.CROSS_PLATFORM:
                opts.authenticatorSelection.authenticatorAttachment = "cross-platform";
                break;
            case ASC.AuthenticatorAttachment.PLATFORM:
                opts.authenticatorSelection.authenticatorAttachment = "platform";
                break;
        }
        switch (auths.getUserVerification()) {
            case ASC.UserVerificationRequirement.DISCOURAGED:
                opts.authenticatorSelection.userVerification = "discouraged";
                break;
            case ASC.UserVerificationRequirement.PREFERRED:
                opts.authenticatorSelection.userVerification = "preferred";
                break;
            case ASC.UserVerificationRequirement.REQUIRED:
                opts.authenticatorSelection.userVerification = "required";
                break;
        }
    }

    opts.challenge = rawOpts.getChallenge_asU8();

    opts.excludeCredentials = credentialsFromProto(rawOpts.getExcludeCredentialsList());

    opts.pubKeyCredParams = [];
    for (const param of rawOpts.getPubKeyCredParamsList()) {
        opts.pubKeyCredParams.push({
            type: "public-key",
            alg: param.getAlg(),
        });
    }

    opts.timeout = 1000 * rawOpts.getTimeout().getSeconds() + rawOpts.getTimeout().getNanos() / 1e6;

    opts.user = {
        displayName: rawOpts.getUser().getDisplayName(),
        name: rawOpts.getUser().getName(),
        id: rawOpts.getUser().getId_asU8(),
    };
    if (rawOpts.getUser().getIcon() != "") {
        (opts.user as any).icon = rawOpts.getUser().getIcon();
    }

    opts.rp = {
        name: rawOpts.getRp().getName(),
        id: rawOpts.getRp().getId(),
    };
    if (rawOpts.getRp().getIcon() != "") {
        (opts.rp as any).icon = rawOpts.getRp().getIcon();
    }

    return opts;
}

function isAttestationResponse(r: AuthenticatorResponse): r is AuthenticatorAttestationResponse {
    return "attestationObject" in r;
}

function isAssertionResponse(r: AuthenticatorResponse): r is AuthenticatorAssertionResponse {
    return "authenticatorData" in r;
}

export function credentialFromJS(input: PublicKeyCredential): C {
    const result = new C();
    result.setId(input.id);
    result.setType(input.type);
    const response = input.response;
    const r = new AR();
    r.setClientDataJson(new Uint8Array(response.clientDataJSON));
    if (isAttestationResponse(response)) {
        const atr = new AAtR();
        atr.setAttestationObject(new Uint8Array(response.attestationObject));
        r.setAttestationResponse(atr);
    } else if (isAssertionResponse(response)) {
        const asr = new AAsR();
        asr.setAuthenticatorData(new Uint8Array(response.authenticatorData));
        asr.setSignature(new Uint8Array(response.signature));
        asr.setUserHandle(new Uint8Array(response.userHandle));
        r.setAssertionResponse(asr);
    }
    result.setResponse(r);
    return result;
}

export function requestOptionsFromProto(input: CRO): PublicKeyCredentialRequestOptions {
    return {
        challenge: input.getChallenge_asU8(),
        timeout: 1000 * input.getTimeout().getSeconds() + input.getTimeout().getNanos() / 1e6,
        allowCredentials: credentialsFromProto(input.getAllowedCredentialsList()),
    };
}
