import 'package:flutter/material.dart';

class SimpleAIDSearchBar extends StatefulWidget {
  final Function(String) onSearch;
  final String hintText;

  const SimpleAIDSearchBar({
    super.key,
    required this.onSearch,
    this.hintText = 'Search AID...',
  });

  @override
  SimpleAIDSearchBarState createState() => SimpleAIDSearchBarState();
}

class SimpleAIDSearchBarState extends State<SimpleAIDSearchBar> {
  final TextEditingController _controller = TextEditingController();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8.0),
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Row(
        children: [
          const Icon(Icons.search, color: Colors.grey),
          const SizedBox(width: 8),
          Expanded(
            child: TextField(
              controller: _controller,
              decoration: InputDecoration(
                hintText: widget.hintText,
                border: InputBorder.none,
              ),
              onSubmitted: (value) => widget.onSearch(value),
            ),
          ),
          TextButton(
            onPressed: () => widget.onSearch(_controller.text),
            child: const Text('搜索'),
          ),
        ],
      ),
    );
  }
}
