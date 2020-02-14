package com.blackh4t;

import sleep.runtime.Scalar;
import javax.crypto.NoSuchPaddingException;
import java.io.File;
import java.io.FileInputStream;
import java.io.ObjectInputStream;
import java.security.KeyPair;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;

public class BeaconTool {

    public static Object ReadObjectFromFile(File f) {
        try {
            if (f.exists()) {
                ObjectInputStream objin = new ObjectInputStream(new FileInputStream(f));
                Scalar result = (Scalar)objin.readObject();
                objin.close();
                return result.objectValue();
            }
        } catch (Exception ex) {
            System.out.println("readObject: " + f + "failed and exception is: " + ex);
        }

        return null;
    }

    public static void Usage() {
        System.out.println("Usage:");
        System.out.println("[*] parse the .beacon_keys to RSA private key and public key in pem format");
        System.out.println("\tBeaconTool -i .beacon_keys -rsa");
        /*
        System.out.println("[*] use the public key from .beacon_keys to decrypt the beacon's online info");
        System.out.println("\tBeaconTool -i .beacon_keys -out online_info.txt");
        System.out.println("[*] use the aes key from the beacon's online info to decrypt transfer data (base64 format)");
        System.out.println("\tBeaconTool -i online_info.txt -aes decrypt");
        System.out.println("[*] use the aes key from the beacon's online info to encrypt transfer data (base64 format)");
        System.out.println("\tBeaconTool -i online_info.txt -aes encrypt");
        */
        System.out.println("[*] compile geacon with the public key from .beacon_keys,which use default c2profile config for communication");
        System.out.println("\tBeaconTool -i .beacon_keys -compile geacon_sourcecode_folder");
    }

    public static void main(String[] args) throws NoSuchPaddingException, NoSuchAlgorithmException {

        int argc = args.length;
        if (argc == 3) {
            File keys = new File(args[1]);
            /*
            File keys = new File(".cobaltstrike.beacon_keys");
            if (!keys.exists()) {
                CommonUtils.writeObject(keys, AsymmetricCrypto.generateKeys());
            }
            */
            KeyPair secret = (KeyPair)ReadObjectFromFile(keys);
            AsymmetricCrypto asymmetricCrypto = new AsymmetricCrypto(secret);
            byte[] publicKey = asymmetricCrypto.exportPublicKey();

            String pemPublicBase64 = Base64.getMimeEncoder().encodeToString(publicKey);
            System.out.println("-----BEGIN PUBLIC KEY-----");
            System.out.println(pemPublicBase64);
            System.out.println("-----END PUBLIC KEY-----");

            byte[] privateKey = asymmetricCrypto.privatekey.getEncoded();
            String pemPrivateBase64 = Base64.getMimeEncoder().encodeToString(privateKey);
            System.out.println("-----BEGIN PRIVATE KEY-----");
            System.out.println(pemPrivateBase64);
            System.out.println("-----END PRIVATE KEY-----");

        } else if (argc == 4) {

        } else {
            Usage();
        }

    }
}
