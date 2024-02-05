import {
  Center,
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  Spinner,
} from "@chakra-ui/react";

interface Props {
  isOpen: boolean;
  onClose?: () => void;
}

export default function ModalSpinner({ isOpen, onClose = () => {} }: Props) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} isCentered>
      <ModalOverlay />

      <ModalContent w="auto">
        <ModalBody>
          <Center>
            <Spinner size="xl" speed="0.6s" />
          </Center>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}
