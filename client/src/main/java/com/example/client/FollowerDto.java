package com.example.client;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class FollowerDto {

        @JsonProperty(value ="my_id")
        @JsonInclude(JsonInclude.Include.NON_NULL)
        private String myId;

        @JsonProperty(value ="friend_id")
        @JsonInclude(JsonInclude.Include.NON_NULL)
        private String friendId;

}




