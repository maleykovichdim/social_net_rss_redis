package com.example.client;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.ToString;

@Data
@ToString
public class PostDto {

        @JsonInclude(JsonInclude.Include.NON_NULL)
        private Long id;

        @JsonProperty(value ="author_id")
        @JsonInclude(JsonInclude.Include.NON_NULL)
        private String authorId;

        @JsonProperty(value ="content")
        @JsonInclude(JsonInclude.Include.NON_NULL)
        private String content;

        @JsonProperty(value = "created_at") //time.Time
        @JsonInclude(JsonInclude.Include.NON_NULL)
        private String createdAt;

}




