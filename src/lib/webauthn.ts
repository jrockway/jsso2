import {
    PublicKeyCredentialCreationOptions as CCO,
    PublicKeyCredentialRequestOptions as CRO,
    PublicKeyCredentialDescriptor as CD,
    AuthenticatorSelectionCriteria as ASC,
    PublicKeyCredential as C,
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
        opts.user.icon = rawOpts.getUser().getIcon();
    }

    opts.rp = {
        name: rawOpts.getRp().getName(),
        id: rawOpts.getRp().getId(),
    };
    if (rawOpts.getRp().getIcon() != "") {
        opts.rp.icon = rawOpts.getRp().getIcon();
    }

    return opts;
}

function fillClientDataJSON(credential: C, response: AuthenticatorResponse): void {
    credential.setClientDataJson(new Uint8Array(response.clientDataJSON));
}

function fillAttestationObject(
    credential: C,
    response: AuthenticatorAttestationResponse | AuthenticatorResponse
): void {
    if ("attestationObject" in response) {
        credential.setAttestationObject(new Uint8Array(response.attestationObject));
    }
}

export function credentialFromJS(input: PublicKeyCredential): C {
    const result = new C();
    result.setId(input.id);
    fillClientDataJSON(result, input.response);
    fillAttestationObject(result, input.response);
    return result;
}

export function requestOptionsFromProto(input: CRO): PublicKeyCredentialRequestOptions {
    return {
        challenge: input.getChallenge_asU8(),
        timeout: 1000 * input.getTimeout().getSeconds() + input.getTimeout().getNanos() / 1e6,
        allowCredentials: credentialsFromProto(input.getAllowedCredentialsList()),
    };
}
