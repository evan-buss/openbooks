import { motion } from "framer-motion";
import React from "react";

export default function AnimatedElement({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <motion.div
      layout
      initial={{ scale: 0.8, opacity: 0 }}
      animate={{ scale: 1, opacity: 1 }}
      exit={{ scale: 0.8, opacity: 0 }}
      transition={{ type: "tween" }}>
      {children}
    </motion.div>
  );
}
