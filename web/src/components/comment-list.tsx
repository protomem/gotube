import { Comment } from "@/domain/entities";
import { ROUTES } from "@/lib/routes";

import NextLink from "next/link";
import {
  Avatar,
  Box,
  Card,
  CardBody,
  Link,
  LinkBox,
  LinkOverlay,
  List,
  Text,
} from "@chakra-ui/react";

type CommentListItemProps = {
  comment: Comment;
};

const CommentListItem = ({ comment }: CommentListItemProps) => {
  return (
    <Card>
      <CardBody
        display="flex"
        flexDirection="row"
        justifyContent="start"
        alignItems="start"
        gap={5}
      >
        <LinkBox>
          <LinkOverlay
            as={NextLink}
            href={`${ROUTES.PROFILE}/${comment.author.nickname}`}
          >
            <Avatar
              name={comment.author.nickname}
              src={comment.author.avatarPath}
            />
          </LinkOverlay>
        </LinkBox>

        <Box>
          <Link
            as={NextLink}
            href={`${ROUTES.PROFILE}/${comment.author.nickname}`}
          >
            <Text fontSize="sm" fontWeight="bold">
              {comment.author.nickname}
            </Text>
          </Link>
          <Text fontSize="md">{comment.content}</Text>
        </Box>
      </CardBody>
    </Card>
  );
};

type CommentListProps = {
  comments: Comment[];
};

const CommentList = ({ comments }: CommentListProps) => {
  return (
    <List spacing={5}>
      {comments.map((comment) => (
        <CommentListItem key={comment.id} comment={comment} />
      ))}
    </List>
  );
};

export default CommentList;
