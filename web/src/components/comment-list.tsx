import { Comment } from "@/domain/entities";

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
import { ROUTES } from "@/lib/routes";

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
          <LinkOverlay as={NextLink} href={`${ROUTES.PROFILE}/${"roman"}`}>
            <Avatar name="Dan Abrahmov" src="https://bit.ly/dan-abramov" />
          </LinkOverlay>
        </LinkBox>

        <Box>
          <Link as={NextLink} href={`${ROUTES.PROFILE}/${"roman"}`}>
            <Text fontSize="sm" fontWeight="bold">
              {"Dan Abramov"}
            </Text>
          </Link>
          <Text fontSize="md">
            {
              "Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis."
            }
          </Text>
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
