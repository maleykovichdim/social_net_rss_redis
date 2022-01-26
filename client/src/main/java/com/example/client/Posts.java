package com.example.client;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestTemplate;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.Random;

public class Posts implements Runnable {

    private Random random =  new Random();
    private RestTemplate restTemplate = new RestTemplate();
    HttpHeaders headers = new HttpHeaders();

    final String URI_POST_TEST = "http://localhost:8080/api/user/postTest";
    String easy = RandomString.digits + "ACEFGHJKLMNPQRUVWXYabcdefhijkprstuvwx";

    public void init(){
        //restTemplate.getMessageConverters().add(new MappingJackson2HttpMessageConverter());
        restTemplate.getMessageConverters().add(0, new MappingJackson2HttpMessageConverter());
//        headers.setContentType(MediaType.APPLICATION_JSON);
//        headers.clearContentHeaders();
        headers.clear();
        headers.add("Content-Type", "application/json;charset=UTF-8");
//        headers.add(HttpHeaders.ACCEPT, MediaType.APPLICATION_JSON_VALUE);
//        headers.add(HttpHeaders.ACCEPT_CHARSET, StandardCharsets.UTF_8.name());
        System.out.println(headers.toString());
    }

    @Override
    public void run() {
        this.init();
        for (int k=0; k< 1000; k++){
            PostDto post = new PostDto();
            RandomString generatedString = new RandomString(23, new SecureRandom(), easy);
            String str = "xxxxxxxxxxxxxxxxxxxxxxxxxxxx";
            System.out.println(generatedString.toString());
            post.setContent(generatedString.toString());
//            post.setContent(str);
            int i = random.nextInt(999)+1;
            post.setAuthorId(String.format("%d", i));
            HttpEntity<PostDto> request = new HttpEntity<>(post, headers);
            StatusResponse addedPost = restTemplate.postForObject( URI_POST_TEST, post, StatusResponse.class);
            System.out.println("Post added : "+addedPost.getStatus()+"   " +post.getAuthorId()+" "+post.getContent());
        }
    }
}
