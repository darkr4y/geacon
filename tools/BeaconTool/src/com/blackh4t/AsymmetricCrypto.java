package com.blackh4t;

import java.io.ByteArrayInputStream;
import java.io.DataInputStream;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.PublicKey;
import javax.crypto.Cipher;
import javax.crypto.NoSuchPaddingException;

public class AsymmetricCrypto {
    public Cipher cipher;
    public PrivateKey privatekey;
    public PublicKey publickey;

    public AsymmetricCrypto(KeyPair pair) throws NoSuchAlgorithmException, NoSuchPaddingException {
        this.privatekey = pair.getPrivate();
        this.publickey = pair.getPublic();
        this.cipher = Cipher.getInstance("RSA/ECB/PKCS1Padding");
    }

    public byte[] exportPublicKey() {
        return this.publickey.getEncoded();
    }

    public byte[] decrypt(byte[] ciphertext) {
        byte[] plaintext = new byte[0];

        try {
            synchronized(this.cipher) {
                this.cipher.init(2, this.privatekey);
                plaintext = this.cipher.doFinal(ciphertext);
            }

            DataInputStream in_handle = new DataInputStream(new ByteArrayInputStream(plaintext));
            int magic = in_handle.readInt();
            if (magic != 48879) {
                System.err.println("Magic number failed :( [RSA decrypt]");
                return new byte[0];
            } else {
                int len = in_handle.readInt();
                if (len > 117) {
                    System.err.println("Length field check failed :( [RSA decrypt]");
                    return new byte[0];
                } else {
                    byte[] result = new byte[len];
                    in_handle.readFully(result, 0, len);
                    return result;
                }
            }
        } catch (Exception ex) {
            System.err.println("RSA decrypt" +  ex);
            return new byte[0];
        }
    }

    public static KeyPair generateKeys() {
        try {
            KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA");
            generator.initialize(1024);
            return generator.generateKeyPair();
        } catch (Exception var1) {
            var1.printStackTrace();
            return null;
        }
    }
}

