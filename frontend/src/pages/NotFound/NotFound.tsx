import { Box, Button, Typography } from "@mui/material";
import { FC, memo } from "react"
import { useNavigate } from "react-router";

const NotFound : FC = () => {
const navigate = useNavigate();

  return (
    <Box
      display="flex" flexDirection="column" alignItems="center"
      justifyContent="center" minHeight="80vh" textAlign="center"
      p={3}
    >
      <Typography variant="h3" gutterBottom>
        404 — Page Not Found
      </Typography>
      <Typography variant="body1" mb={2}>
        Sorry, the page you are looking for does not exist.
      </Typography>
      <Button variant="contained" onClick={() => navigate('/')}>
        Go to Home
      </Button>
    </Box>
  );
}

export default memo(NotFound)