import {
  IconButton,
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import { FaUpload } from "react-icons/fa6";
import UploadVideoForm from "@/components/upload-video-form";

export default function UploadVideoMenu() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton aria-label="Create Video" onClick={onOpen}>
        <FaUpload />
      </IconButton>

      <Modal isOpen={isOpen} onClose={onClose} isCentered size="4xl">
        <ModalOverlay />
        <ModalContent>
          <ModalBody p="0" my="6" display="flex" justifyContent="center">
            <UploadVideoForm afterClose={onClose} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}
